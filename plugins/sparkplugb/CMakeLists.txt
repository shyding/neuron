# CMakeLists.txt for sparkplugb
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/sparkplugb.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: sparkplugb.so
set(plugin_sparkplugb_SOURCES sparkplugb.c sparkplugb.pb.c sparkplugb_plugin.c util.c)
add_library(plugin-sparkplugb SHARED ${plugin_sparkplugb_SOURCES})
target_include_directories(plugin-sparkplugb PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-sparkplugb neuron-base)

