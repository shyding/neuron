# CMakeLists.txt for secs_gem
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/hsms.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: hsms.so
set(plugin_hsms_SOURCES hsms.c hsms_plugin.c hsms_stack.c secs2.c secs_gem_point.c)
add_library(plugin-hsms SHARED ${plugin_hsms_SOURCES})
target_include_directories(plugin-hsms PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-hsms neuron-base)

