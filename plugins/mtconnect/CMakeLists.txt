# CMakeLists.txt for mtconnect
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/mtc.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: mtconnect.so
set(plugin_mtconnect_SOURCES mtc_plugin.c mtc_point.c mtc_stack.c)
add_library(plugin-mtconnect SHARED ${plugin_mtconnect_SOURCES})
target_include_directories(plugin-mtconnect PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-mtconnect neuron-base)

