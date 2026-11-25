# CMakeLists.txt for iec103
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/iec103.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: iec103.so
set(plugin_iec103_SOURCES iec103.c iec103_plugin.c iec103_stack.c iec103_tag.c)
add_library(plugin-iec103 SHARED ${plugin_iec103_SOURCES})
target_include_directories(plugin-iec103 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-iec103 neuron-base)

