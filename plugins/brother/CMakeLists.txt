# CMakeLists.txt for brother
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/brother.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: brother.so
set(plugin_brother_SOURCES brother.c brother_plugin.c brother_stack.c brother_tag.c)
add_library(plugin-brother SHARED ${plugin_brother_SOURCES})
target_include_directories(plugin-brother PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-brother neuron-base)

