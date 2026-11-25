# CMakeLists.txt for codesys3
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/codesys3.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: codesys3.so
set(plugin_codesys3_SOURCES codesys3.c codesys3_plugin.c codesys3_point.c codesys3_stack.c)
add_library(plugin-codesys3 SHARED ${plugin_codesys3_SOURCES})
target_include_directories(plugin-codesys3 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-codesys3 neuron-base)

