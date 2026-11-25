# CMakeLists.txt for heidenhain_cnc
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/heidenhain_cnc.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: heidenhain_cnc.so
set(plugin_heidenhain_cnc_SOURCES heidenhain_cnc_plugin.c lsv2.c lsv2_point.c lsv2_stack.c)
add_library(plugin-heidenhain_cnc SHARED ${plugin_heidenhain_cnc_SOURCES})
target_include_directories(plugin-heidenhain_cnc PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-heidenhain_cnc neuron-base)

