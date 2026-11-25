# CMakeLists.txt for iec102
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/iec102.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: iec102.so
set(plugin_iec102_SOURCES iec102.c iec102_plugin.c iec102_stack.c iec102_tag.c)
add_library(plugin-iec102 SHARED ${plugin_iec102_SOURCES})
target_include_directories(plugin-iec102 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-iec102 neuron-base)

