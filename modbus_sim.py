import socket
import struct
import threading
import time

class ModbusMultiSlaveSim:
    def __init__(self, host='0.0.0.0', port=10502):
        self.host = host
        self.port = port
        self.running = True
        self.sock = None

        # Registers state per slave unit_id
        # Slave 1: sub-pump-01 -> [pressure, flow, vibration, runState]
        # Slave 2: sub-pump-02 -> [pressure, flow, vibration, runState]
        # Slave 3: sub-valve-01 -> [opening, status]
        self.slave_registers = {
            1: [165, 4520, 42, 1],
            2: [158, 4380, 38, 1],
            3: [75, 1, 0, 0],
            0: [100, 200, 300, 400]
        }
        self.lock = threading.Lock()
        self.t = threading.Thread(target=self._run, daemon=True)

    def start(self):
        self.t.start()
        time.sleep(0.5)

    def _run(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self.sock.bind((self.host, self.port))
        self.sock.listen(10)
        print(f"Modbus Multi-Slave Sim listening on {self.host}:{self.port}")

        while self.running:
            try:
                conn, addr = self.sock.accept()
                threading.Thread(target=self._handle_client, args=(conn,), daemon=True).start()
            except:
                break

    def _handle_client(self, conn):
        with conn:
            while self.running:
                try:
                    data = conn.recv(1024)
                except:
                    break
                if not data or len(data) < 12:
                    break

                trans_id, proto_id, length, unit_id, fc = struct.unpack(">HHHBB", data[:8])

                if fc == 3 or fc == 4:  # Read Holding / Input Registers
                    addr, qty = struct.unpack(">HH", data[8:12])
                    byte_cnt = qty * 2
                    resp_len = 3 + byte_cnt
                    resp_hdr = struct.pack(">HHHBB", trans_id, 0, resp_len, unit_id, fc)
                    resp_body = struct.pack("B", byte_cnt)

                    with self.lock:
                        regs = self.slave_registers.get(unit_id, [100, 200, 300, 400])
                        now_int = int(time.time())

                        # Dynamic variation based on slave
                        if unit_id == 1:
                            val_list = [
                                (now_int % 40) + 160,       # 1.60 ~ 2.00 MPa
                                (now_int % 100) + 4500,     # 45.00 ~ 46.00 m3/h
                                (now_int % 15) + 38,        # 0.38 ~ 0.53 mm/s
                                regs[3]                     # runState
                            ]
                        elif unit_id == 2:
                            val_list = [
                                (now_int % 35) + 155,       # 1.55 ~ 1.90 MPa
                                (now_int % 90) + 4300,      # 43.00 ~ 43.90 m3/h
                                (now_int % 12) + 35,        # 0.35 ~ 0.47 mm/s
                                regs[3]                     # runState
                            ]
                        elif unit_id == 3:
                            val_list = [
                                regs[0],                    # opening (e.g. 75%)
                                regs[1],                    # status (1 = open)
                                0,
                                0
                            ]
                        else:
                            val_list = regs

                    for i in range(qty):
                        reg_idx = (addr + i) % len(val_list)
                        val = val_list[reg_idx]
                        resp_body += struct.pack(">H", val)

                    try:
                        conn.sendall(resp_hdr + resp_body)
                    except:
                        break

                elif fc == 6:  # Write Single Register
                    addr, val = struct.unpack(">HH", data[8:12])
                    with self.lock:
                        if unit_id not in self.slave_registers:
                            self.slave_registers[unit_id] = [0, 0, 0, 0]
                        if addr < len(self.slave_registers[unit_id]):
                            self.slave_registers[unit_id][addr] = val
                    try:
                        conn.sendall(data[:12])
                    except:
                        break

                elif fc == 16:  # Write Multiple Registers
                    addr, qty, byte_cnt = struct.unpack(">HHB", data[8:13])
                    vals = struct.unpack(f">{qty}H", data[13:13 + byte_cnt])
                    with self.lock:
                        if unit_id not in self.slave_registers:
                            self.slave_registers[unit_id] = [0] * max(4, addr + qty)
                        while len(self.slave_registers[unit_id]) < addr + qty:
                            self.slave_registers[unit_id].append(0)
                        for i, v in enumerate(vals):
                            self.slave_registers[unit_id][addr + i] = v
                    resp_hdr = struct.pack(">HHHBBHH", trans_id, 0, 6, unit_id, fc, addr, qty)
                    try:
                        conn.sendall(resp_hdr)
                    except:
                        break
                else:
                    break

if __name__ == '__main__':
    sim = ModbusMultiSlaveSim()
    sim.start()
    print("Modbus Multi-Slave Sim running on 0.0.0.0:10502 (Slaves 1, 2, 3)")
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        pass
