# CMakeLists.txt for hostlink
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/hostlink-cmode.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: hostlink-cmode.so
set(plugin_hostlink_cmode_SOURCES hostlink.c hostlink_cmode.c hostlink_point.c hostlink_req.c hostlink_stack.c)
add_library(plugin-hostlink-cmode SHARED ${plugin_hostlink_cmode_SOURCES})
target_include_directories(plugin-hostlink-cmode PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-hostlink-cmode neuron-base)

