# CMakeLists.txt for knx
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/knx.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: knx.so
set(plugin_knx_SOURCES knx_cemi.c knx_client.c knx_frame.c knx_plugin.c)
add_library(plugin-knx SHARED ${plugin_knx_SOURCES})
target_include_directories(plugin-knx PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-knx neuron-base)

