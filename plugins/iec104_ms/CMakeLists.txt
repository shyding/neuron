# CMakeLists.txt for iec104_ms
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/iec104_ms.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

set(plugin_iec104_ms_SOURCES 
    iec104_ms_plugin.c 
    iec104_ms_point.c 
    iec104_ms_stack.c
    ../iec104/iec104.c
)
add_library(plugin-iec104-ms SHARED ${plugin_iec104_ms_SOURCES})
target_include_directories(plugin-iec104-ms PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_SOURCE_DIR}/plugins/iec104
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-iec104-ms neuron-base)
