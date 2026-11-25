set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")

set(MODBUS_SRC modbus.c modbus_point.c modbus_req.c modbus_stack.c)

set(CMAKE_BUILD_RPATH ./)
file(COPY ${CMAKE_SOURCE_DIR}/plugins/modbus/modbus-tcp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_SOURCE_DIR}/plugins/modbus/modbus-rtu.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_SOURCE_DIR}/plugins/modbus/modbus-ascii.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_SOURCE_DIR}/plugins/modbus/modbus-qh-tcp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_SOURCE_DIR}/plugins/modbus/hollysys-modbus-rtu.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_SOURCE_DIR}/plugins/modbus/hollysys-modbus-tcp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_SOURCE_DIR}/plugins/modbus/inovance-modbus-tcp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_SOURCE_DIR}/plugins/modbus/xinje-modbus-rtu.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# modbus rtu plugin
set(MODBUS_RTU_PLUGIN plugin-modbus-rtu)
set(MODBUS_RTU_PLUGIN_SOURCES  modbus_rtu.c
                               ${MODBUS_SRC})
add_library(${MODBUS_RTU_PLUGIN} SHARED)
target_include_directories(${MODBUS_RTU_PLUGIN} PRIVATE ${CMAKE_SOURCE_DIR}/include/neuron ${CMAKE_SOURCE_DIR}/plugins/include)
target_sources(${MODBUS_RTU_PLUGIN} PRIVATE ${MODBUS_RTU_PLUGIN_SOURCES})
target_link_libraries(${MODBUS_RTU_PLUGIN} neuron-base)

# modbus tcp plugin
set(MODBUS_TCP_PLUS_PLUGIN plugin-modbus-tcp)
set(MODBUS_TCP_PLUS_PLUGIN_SOURCES  modbus_tcp.c
                                    ${MODBUS_SRC})
add_library(${MODBUS_TCP_PLUS_PLUGIN} SHARED)
target_include_directories(${MODBUS_TCP_PLUS_PLUGIN} PRIVATE ${CMAKE_SOURCE_DIR}/include/neuron ${CMAKE_SOURCE_DIR}/plugins/include)
target_sources(${MODBUS_TCP_PLUS_PLUGIN} PRIVATE ${MODBUS_TCP_PLUS_PLUGIN_SOURCES})
target_link_libraries(${MODBUS_TCP_PLUS_PLUGIN} neuron-base)

# modbus ascii plugin
set(MODBUS_ASCII_PLUGIN plugin-modbus-ascii)
add_library(${MODBUS_ASCII_PLUGIN} SHARED modbus_ascii.c ${MODBUS_SRC} hollysys_modbus_point.c xinje_modbus_point.c)
target_include_directories(${MODBUS_ASCII_PLUGIN} PRIVATE ${CMAKE_SOURCE_DIR}/include/neuron ${CMAKE_SOURCE_DIR}/plugins/include)
target_link_libraries(${MODBUS_ASCII_PLUGIN} neuron-base)

# modbus qh tcp plugin
set(MODBUS_QH_TCP_PLUGIN plugin-modbus-qh-tcp)
add_library(${MODBUS_QH_TCP_PLUGIN} SHARED modbus_qh_tcp.c ${MODBUS_SRC} hollysys_modbus_point.c xinje_modbus_point.c)
target_include_directories(${MODBUS_QH_TCP_PLUGIN} PRIVATE ${CMAKE_SOURCE_DIR}/include/neuron ${CMAKE_SOURCE_DIR}/plugins/include)
target_link_libraries(${MODBUS_QH_TCP_PLUGIN} neuron-base)

# hollysys modbus rtu plugin
set(HOLLYSYS_RTU_PLUGIN plugin-hollysys-modbus-rtu)
add_library(${HOLLYSYS_RTU_PLUGIN} SHARED hollysys_modbus_rtu.c hollysys_modbus_point.c ${MODBUS_SRC} xinje_modbus_point.c)
target_include_directories(${HOLLYSYS_RTU_PLUGIN} PRIVATE ${CMAKE_SOURCE_DIR}/include/neuron ${CMAKE_SOURCE_DIR}/plugins/include)
target_link_libraries(${HOLLYSYS_RTU_PLUGIN} neuron-base)

# hollysys modbus tcp plugin
set(HOLLYSYS_TCP_PLUGIN plugin-hollysys-modbus-tcp)
add_library(${HOLLYSYS_TCP_PLUGIN} SHARED hollysys_modbus_tcp.c hollysys_modbus_point.c ${MODBUS_SRC} xinje_modbus_point.c)
target_include_directories(${HOLLYSYS_TCP_PLUGIN} PRIVATE ${CMAKE_SOURCE_DIR}/include/neuron ${CMAKE_SOURCE_DIR}/plugins/include)
target_link_libraries(${HOLLYSYS_TCP_PLUGIN} neuron-base)

# inovance modbus tcp plugin
set(INOVANCE_TCP_PLUGIN plugin-inovance-modbus-tcp)
add_library(${INOVANCE_TCP_PLUGIN} SHARED inovance_modbus_tcp.c hollysys_modbus_point.c ${MODBUS_SRC} xinje_modbus_point.c)
target_include_directories(${INOVANCE_TCP_PLUGIN} PRIVATE ${CMAKE_SOURCE_DIR}/include/neuron ${CMAKE_SOURCE_DIR}/plugins/include)
target_link_libraries(${INOVANCE_TCP_PLUGIN} neuron-base)

# xinje modbus rtu plugin
set(XINJE_RTU_PLUGIN plugin-xinje-modbus-rtu)
add_library(${XINJE_RTU_PLUGIN} SHARED xinje_modbus_rtu.c xinje_modbus_point.c hollysys_modbus_point.c ${MODBUS_SRC})
target_include_directories(${XINJE_RTU_PLUGIN} PRIVATE ${CMAKE_SOURCE_DIR}/include/neuron ${CMAKE_SOURCE_DIR}/plugins/include)
target_link_libraries(${XINJE_RTU_PLUGIN} neuron-base)
