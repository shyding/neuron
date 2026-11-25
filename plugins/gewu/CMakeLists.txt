# CMakeLists.txt for gewu
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/gewu.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: gewu.so
set(plugin_gewu_SOURCES gewu_config.c gewu_handle.c gewu_plugin.c)
add_library(plugin-gewu SHARED ${plugin_gewu_SOURCES})
target_include_directories(plugin-gewu PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-gewu neuron-base)

