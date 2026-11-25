# CMakeLists.txt for EtherNet_IP
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/EtherNet-IP.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/ethernet-ip-1400.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/ethernet-ip-55.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: EtherNet-IP-1400.so
set(plugin_EtherNet_IP_1400_SOURCES cip.c enet.c enet_1400_plugin.c enet_1400_stack.c enet_1400_tag.c)
add_library(plugin-EtherNet-IP-1400 SHARED ${plugin_EtherNet_IP_1400_SOURCES})
target_include_directories(plugin-EtherNet-IP-1400 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-EtherNet-IP-1400 neuron-base)

# Plugin: EtherNet-IP-5500.so
set(plugin_EtherNet_IP_5500_SOURCES cip.c enet.c enet_55_plugin.c enet_55_stack.c enet_55_tag.c)
add_library(plugin-EtherNet-IP-5500 SHARED ${plugin_EtherNet_IP_5500_SOURCES})
target_include_directories(plugin-EtherNet-IP-5500 PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-EtherNet-IP-5500 neuron-base)

# Plugin: EtherNet-IP.so
set(plugin_EtherNet_IP_SOURCES cip.c enet.c enet_plugin.c enet_stack.c enet_tag.c)
add_library(plugin-EtherNet-IP SHARED ${plugin_EtherNet_IP_SOURCES})
target_include_directories(plugin-EtherNet-IP PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-EtherNet-IP neuron-base)

