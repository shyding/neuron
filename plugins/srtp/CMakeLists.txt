# CMakeLists.txt for srtp
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/srtp.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: srtp.so
set(plugin_srtp_SOURCES srtp.c srtp_plugin.c srtp_point.c srtp_stack.c)
add_library(plugin-srtp SHARED ${plugin_srtp_SOURCES})
target_include_directories(plugin-srtp PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-srtp neuron-base)

