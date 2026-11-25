# CMakeLists.txt for kuka
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/kuka.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: kuka.so
set(plugin_kuka_SOURCES kuka_plugin.c kuka_point.c kuka_stack.c)
add_library(plugin-kuka SHARED ${plugin_kuka_SOURCES})
target_include_directories(plugin-kuka PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-kuka neuron-base)

