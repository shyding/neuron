# CMakeLists.txt for ads
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/ads.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: ads.so
set(plugin_ads_SOURCES ads.c ads_client.c ads_plugin.c ads_tag.c xio/ads_xio_posix.c)
add_library(plugin-ads SHARED ${plugin_ads_SOURCES})
target_include_directories(plugin-ads PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-ads neuron-base)

