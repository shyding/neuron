set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/cctl.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

set(PLUGIN_NAME plugin-cctl)
set(PLUGIN_SOURCES
    cctl_plugin.c
    cctl_handle.c
    cctl_config.c
    cctl_tcp.c
)

add_library(${PLUGIN_NAME} SHARED)
target_include_directories(${PLUGIN_NAME} PRIVATE
    ${CMAKE_SOURCE_DIR}/include/neuron
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_sources(${PLUGIN_NAME} PRIVATE ${PLUGIN_SOURCES})
target_link_libraries(${PLUGIN_NAME} neuron-base jansson)
