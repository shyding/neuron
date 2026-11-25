# CMakeLists.txt for fx
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/fx.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: fx.so
set(plugin_fx_SOURCES fx_plugin.c fx_point.c fx_protocol.c)
add_library(plugin-fx SHARED ${plugin_fx_SOURCES})
target_include_directories(plugin-fx PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-fx neuron-base)

