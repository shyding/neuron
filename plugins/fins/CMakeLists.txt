# CMakeLists.txt for fins
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/fins-tcp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/fins-udp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: fins-tcp.so
set(plugin_fins_tcp_SOURCES fins.c fins_plugin.c fins_point.c fins_stack.c)
add_library(plugin-fins-tcp SHARED ${plugin_fins_tcp_SOURCES})
target_include_directories(plugin-fins-tcp PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-fins-tcp neuron-base)

# Plugin: fins-udp.so
set(plugin_fins_udp_SOURCES fins.c fins_point.c fins_stack.c fins_udp.c)
add_library(plugin-fins-udp SHARED ${plugin_fins_udp_SOURCES})
target_include_directories(plugin-fins-udp PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-fins-udp neuron-base)

