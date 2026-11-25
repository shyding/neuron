# CMakeLists.txt for opcua
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/opcua.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: opcua.so
set(plugin_opcua_SOURCES opcua_browse_cache.c opcua_client_common.c opcua_client_subscribe.c opcua_condition.c opcua_handle.c opcua_plugin.c open62541_client.c)
add_library(plugin-opcua SHARED ${plugin_opcua_SOURCES})
target_include_directories(plugin-opcua PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-opcua neuron-base)

