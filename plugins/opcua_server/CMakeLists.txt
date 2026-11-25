# CMakeLists.txt for opcua_server
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/opcuaserver.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: opcua-server.so
set(plugin_opcua_server_SOURCES cert_util.c my_ua_accesscontrol.c opcua_server_handle.c opcua_server_node_manager.c opcua_server_plugin.c)
add_library(plugin-opcua-server SHARED ${plugin_opcua_server_SOURCES})
target_include_directories(plugin-opcua-server PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-opcua-server neuron-base)

