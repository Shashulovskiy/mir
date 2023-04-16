#include <vector>
#include <cstdint>

std::vector<uint8_t> encode(uint32_t k, std::vector<uint8_t> &data);
std::vector<uint8_t> decode(uint32_t k, std::vector<uint8_t> &data, std::vector<int> &missing);