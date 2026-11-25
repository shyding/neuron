# CMakeLists.txt for hj
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/hj.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: HJ212.so
set(plugin_HJ212_SOURCES hj.c hj_plugin.c hj_stack.c hj_tag.c)
add_library(plugin-HJ212 SHARED ${plugin_HJ212_SOURCES})
target_include_directories(plugin-HJ212 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-HJ212 neuron-base)

