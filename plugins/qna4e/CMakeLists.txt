# CMakeLists.txt for qna4e
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/qna4e.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: qna4e.so
set(plugin_qna4e_SOURCES qna4e_plugin.c qna4e_point.c qna4e_protocol.c)
add_library(plugin-qna4e SHARED ${plugin_qna4e_SOURCES})
target_include_directories(plugin-qna4e PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-qna4e neuron-base)

