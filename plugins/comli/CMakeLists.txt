# CMakeLists.txt for comli
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/comli.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: comli.so
set(plugin_comli_SOURCES comli.c comli_plugin.c comli_point.c)
add_library(plugin-comli SHARED ${plugin_comli_SOURCES})
target_include_directories(plugin-comli PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-comli neuron-base)

