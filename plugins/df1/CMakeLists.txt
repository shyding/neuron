# CMakeLists.txt for df1
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/df1.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: df1.so
set(plugin_df1_SOURCES df1.c df1_plugin.c df1_point.c df1_req.c df1_stack.c)
add_library(plugin-df1 SHARED ${plugin_df1_SOURCES})
target_include_directories(plugin-df1 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-df1 neuron-base)

