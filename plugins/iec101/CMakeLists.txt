# CMakeLists.txt for iec101
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/iec101.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: iec101.so
set(plugin_iec101_SOURCES iec101.c iec101_plugin.c iec101_stack.c)
add_library(plugin-iec101 SHARED ${plugin_iec101_SOURCES})
target_include_directories(plugin-iec101 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-iec101 neuron-base)

