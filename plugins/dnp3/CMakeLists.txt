# CMakeLists.txt for dnp3
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/dnp3.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: dnp3.so
set(plugin_dnp3_SOURCES dnp3.c dnp3_plugin.c dnp3_point.c dnp3_stack.c)
add_library(plugin-dnp3 SHARED ${plugin_dnp3_SOURCES})
target_include_directories(plugin-dnp3 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-dnp3 neuron-base)

