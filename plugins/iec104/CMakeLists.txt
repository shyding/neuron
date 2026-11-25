# CMakeLists.txt for iec104
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/iec104_ms.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/iec104_standard.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: iec104-ms.so
set(plugin_iec104_ms_SOURCES iec104.c)
add_library(plugin-iec104-ms SHARED ${plugin_iec104_ms_SOURCES})
target_include_directories(plugin-iec104-ms PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-iec104-ms neuron-base)

# Plugin: iec104-standard.so
set(plugin_iec104_standard_SOURCES iec104.c iec104_stack.c iec104_standard_plugin.c)
add_library(plugin-iec104-standard SHARED ${plugin_iec104_standard_SOURCES})
target_include_directories(plugin-iec104-standard PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-iec104-standard neuron-base)

