# CMakeLists.txt for knd
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/knd.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: knd.so
set(plugin_knd_SOURCES knd_handle.c knd_plugin.c knd_point.c)
add_library(plugin-knd SHARED ${plugin_knd_SOURCES})
target_include_directories(plugin-knd PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-knd neuron-base)

