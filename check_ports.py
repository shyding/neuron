import socket

ports = [7000, 8085, 8086, 8088]
for p in ports:
    try:
        s = socket.socket()
        s.settimeout(1.0)
        s.connect(('127.0.0.1', p))
        s.close()
        print(f"Port {p}: LISTENING")
    except Exception as e:
        print(f"Port {p}: CLOSED ({e})")
