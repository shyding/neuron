# CMakeLists.txt for cip
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/ab5000.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/omron-cip.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: ab5000.so
set(plugin_ab5000_SOURCES ab5000_plugin.c ab_stack.c cip.c enip.c tag.c)
add_library(plugin-ab5000 SHARED ${plugin_ab5000_SOURCES})
target_include_directories(plugin-ab5000 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-ab5000 neuron-base)

# Plugin: omron-cip.so
set(plugin_omron_cip_SOURCES cip.c enip.c omron_plugin.c stack.c tag.c)
add_library(plugin-omron-cip SHARED ${plugin_omron_cip_SOURCES})
target_include_directories(plugin-omron-cip PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-omron-cip neuron-base)

