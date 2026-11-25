# CMakeLists.txt for mewtocol
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/mewtocol.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: mewtocol.so
set(plugin_mewtocol_SOURCES mewtocol_plugin.c mewtocol_point.c mewtocol_protocol.c)
add_library(plugin-mewtocol SHARED ${plugin_mewtocol_SOURCES})
target_include_directories(plugin-mewtocol PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-mewtocol neuron-base)

