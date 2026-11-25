# CMakeLists.txt for websocket
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/websocket.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: websocket.so
set(plugin_websocket_SOURCES websocket_client.c websocket_config.c websocket_handle.c websocket_plugin.c)
add_library(plugin-websocket SHARED ${plugin_websocket_SOURCES})
target_include_directories(plugin-websocket PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-websocket neuron-base)

