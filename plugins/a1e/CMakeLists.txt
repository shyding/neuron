# CMakeLists.txt for a1e
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/a1e.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: a1e.so
set(plugin_a1e_SOURCES a1e_plugin.c a1e_point.c a1e_protocol.c)
add_library(plugin-a1e SHARED ${plugin_a1e_SOURCES})
target_include_directories(plugin-a1e PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-a1e neuron-base)

