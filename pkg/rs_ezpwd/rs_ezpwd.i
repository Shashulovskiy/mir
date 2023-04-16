%module rs_ezpwd

%{
#include <vector>
#include "rs_ezpwd.hpp"
%}

%include <std_vector.i>

%include <typemaps.i>
%include "stdint.i"

// Wrap the vector of vectors
%template(VecUInt8) std::vector<unsigned char>;
%template(VecInt) std::vector<int>;
%template(VecVecUInt8) std::vector<std::vector<unsigned char>>;

// Wrap coding
std::vector<uint8_t> encode(uint32_t k, std::vector<uint8_t> &data);
std::vector<uint8_t> decode(uint32_t k, std::vector<uint8_t> &data, std::vector<int> &missing);