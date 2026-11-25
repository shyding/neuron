# CMakeLists.txt for bacnet
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/bacnet.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/bacnet_scan.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: bacnet-scan.so
set(plugin_bacnet_scan_SOURCES scanning/protocol/bacnet_apdu.c scanning/protocol/bacnet_bvll.c scanning/protocol/bacnet_confirmed.c scanning/protocol/bacnet_error.c scanning/protocol/bacnet_npdu.c scanning/protocol/bacnet_object.c scanning/protocol/bacnet_property.c scanning/protocol/bacnet_read.c scanning/protocol/bacnet_service.c scanning/protocol/bacnet_units.c scanning/scan_cache.c scanning/scan_config.c scanning/scan_device.c scanning/scan_discovery.c scanning/scan_export.c scanning/scan_object.c scanning/scan_plugin.c scanning/scan_point.c scanning/scan_router.c scanning/scan_schedule.c scanning/scan_snapshot.c scanning/scan_transaction.c scanning/scan_transport.c)
add_library(plugin-bacnet-scan SHARED ${plugin_bacnet_scan_SOURCES})
target_include_directories(plugin-bacnet-scan PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-bacnet-scan neuron-base)

# Plugin: bacnet.so
set(plugin_bacnet_SOURCES bacapp.c bacnet.c bacnet_plugin.c bacnet_point.c bacnet_stack.c)
add_library(plugin-bacnet SHARED ${plugin_bacnet_SOURCES})
target_include_directories(plugin-bacnet PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-bacnet neuron-base)

