# CMakeLists.txt for nona11
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/nona11.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: nona11.so
set(plugin_nona11_SOURCES nona11.c nona11_plugin.c nona11_stack.c nona11_tag.c)
add_library(plugin-nona11 SHARED ${plugin_nona11_SOURCES})
target_include_directories(plugin-nona11 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-nona11 neuron-base)

