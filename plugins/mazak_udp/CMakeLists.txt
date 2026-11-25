# CMakeLists.txt for mazak_udp
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/mazak-udp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: mazak-udp.so
set(plugin_mazak_udp_SOURCES mazak_udp_config.c mazak_udp_parser.c mazak_udp_plugin.c mazak_udp_tag.c)
add_library(plugin-mazak-udp SHARED ${plugin_mazak_udp_SOURCES})
target_include_directories(plugin-mazak-udp PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-mazak-udp neuron-base)

