# CMakeLists.txt for iec61850
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/iec61850.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: iec61850.so
set(plugin_iec61850_SOURCES asn.c iec61850_plugin.c iec61850_protocol.c iec61850_stack.c tpkt.c)
add_library(plugin-iec61850 SHARED ${plugin_iec61850_SOURCES})
target_include_directories(plugin-iec61850 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-iec61850 neuron-base)

