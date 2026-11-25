# CMakeLists.txt for dlt645
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/dlt645-1997.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/dlt645.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: dlt645-1997.so
set(plugin_dlt645_1997_SOURCES dlt645.c dlt645_97.c dlt645_point.c dlt645_req.c dlt645_stack.c)
add_library(plugin-dlt645-1997 SHARED ${plugin_dlt645_1997_SOURCES})
target_include_directories(plugin-dlt645-1997 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-dlt645-1997 neuron-base)

# Plugin: dlt645-2007.so
set(plugin_dlt645_2007_SOURCES dlt645.c dlt645_plugin.c dlt645_point.c dlt645_req.c dlt645_stack.c)
add_library(plugin-dlt645-2007 SHARED ${plugin_dlt645_2007_SOURCES})
target_include_directories(plugin-dlt645-2007 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-dlt645-2007 neuron-base)

