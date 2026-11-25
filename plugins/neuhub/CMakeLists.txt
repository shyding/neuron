# CMakeLists.txt for neuhub
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/neuhub.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: neuhub.so
set(plugin_neuhub_SOURCES neuhub.c neuhub_plugin.c neuhub_point.c neuhub_stack.c)
add_library(plugin-neuhub SHARED ${plugin_neuhub_SOURCES})
target_include_directories(plugin-neuhub PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-neuhub neuron-base)

