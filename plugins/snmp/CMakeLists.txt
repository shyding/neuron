# CMakeLists.txt for snmp
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/snmp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: snmp.so
set(plugin_snmp_SOURCES asn.c snmp.c snmp_plugin.c snmp_stack.c snmp_tag.c)
add_library(plugin-snmp SHARED ${plugin_snmp_SOURCES})
target_include_directories(plugin-snmp PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-snmp neuron-base)

