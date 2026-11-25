# CMakeLists.txt for qna3e
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/qna3e.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: qna3e.so
set(plugin_qna3e_SOURCES qna3e_plugin.c qna3e_point.c qna3e_protocol.c)
add_library(plugin-qna3e SHARED ${plugin_qna3e_SOURCES})
target_include_directories(plugin-qna3e PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-qna3e neuron-base)

