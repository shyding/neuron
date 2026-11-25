# CMakeLists.txt for mpi
set(LIBRARY_OUTPUT_PATH "${CMAKE_BINARY_DIR}/plugins")
set(CMAKE_BUILD_RPATH ./)

# Copy schemas
file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/mpi.json DESTINATION ${CMAKE_BINARY_DIR}/plugins/schema/)

# Plugin: mpi.so
set(plugin_mpi_SOURCES mpi_plugin.c mpi_point.c mpi_protocol.c)
add_library(plugin-mpi SHARED ${plugin_mpi_SOURCES})
target_include_directories(plugin-mpi PRIVATE 
    ${CMAKE_SOURCE_DIR}/include/neuron 
    ${CMAKE_SOURCE_DIR}/plugins/include
    ${CMAKE_CURRENT_SOURCE_DIR}
)
target_link_libraries(plugin-mpi neuron-base)

