# CMakeLists.txt for s5fw
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/s5fetch-write.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: s5fetch-write.so
set(plugin_s5fetch_write_SOURCES s5fw_plugin.c s5fw_point.c s5fw_protocol.c)
add_library(plugin-s5fetch-write SHARED ${plugin_s5fetch_write_SOURCES})
target_include_directories(plugin-s5fetch-write PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-s5fetch-write neuron-base)

