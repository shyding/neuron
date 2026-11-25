# CMakeLists.txt for s7comm
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/s7comm.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/s7comm_for_300.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: s7comm-for-300.so
set(plugin_s7comm_for_300_SOURCES cotp.c s7comm.c s7comm_300_plugin.c s7comm_point.c s7comm_stack.c tpkt.c)
add_library(plugin-s7comm-for-300 SHARED ${plugin_s7comm_for_300_SOURCES})
target_include_directories(plugin-s7comm-for-300 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-s7comm-for-300 neuron-base)

# Plugin: s7comm.so
set(plugin_s7comm_SOURCES cotp.c s7comm.c s7comm_plugin.c s7comm_point.c s7comm_stack.c tpkt.c)
add_library(plugin-s7comm SHARED ${plugin_s7comm_SOURCES})
target_include_directories(plugin-s7comm PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-s7comm neuron-base)

