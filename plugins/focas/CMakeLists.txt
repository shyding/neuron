# CMakeLists.txt for focas
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/focas.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: focas.so
set(plugin_focas_SOURCES focas_plugin.c focas_stack.c focas_tag.c)
add_library(plugin-focas SHARED ${plugin_focas_SOURCES})
target_include_directories(plugin-focas PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-focas neuron-base)

