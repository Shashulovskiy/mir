//
// Created by Шашуловский Артем Владимирович on 08.04.2023.
//

#include "ezpwd-reed-solomon/c++/ezpwd/rs"
#include "rs_ezpwd.hpp"
#include <vector>

template<uint32_t f>
std::vector<uint8_t> encode(ezpwd::RS<255, 255 - 2 * f> &coder, std::vector<uint8_t> &data) {
    coder.encode(data);

    return data;
}

template<uint32_t f>
std::vector<uint8_t> decode(ezpwd::RS<255, 255 - 2 * f> &coder, std::vector<uint8_t> &data, std::vector<int> &missing) {
    int corrected = coder.decode(data, missing);
    if (corrected >= 0) {
        data.resize(data.size() - coder.nroots());
        return data;
    } else {
        return {};
    }
}

ezpwd::RS<255, 255 - 2 * 1> coder1;
ezpwd::RS<255, 255 - 2 * 2> coder2;
ezpwd::RS<255, 255 - 2 * 3> coder3;
ezpwd::RS<255, 255 - 2 * 4> coder4;
ezpwd::RS<255, 255 - 2 * 5> coder5;
ezpwd::RS<255, 255 - 2 * 6> coder6;
ezpwd::RS<255, 255 - 2 * 7> coder7;
ezpwd::RS<255, 255 - 2 * 8> coder8;
ezpwd::RS<255, 255 - 2 * 9> coder9;
ezpwd::RS<255, 255 - 2 * 10> coder10;
ezpwd::RS<255, 255 - 2 * 11> coder11;
ezpwd::RS<255, 255 - 2 * 12> coder12;
ezpwd::RS<255, 255 - 2 * 13> coder13;
ezpwd::RS<255, 255 - 2 * 14> coder14;
ezpwd::RS<255, 255 - 2 * 15> coder15;
ezpwd::RS<255, 255 - 2 * 16> coder16;
ezpwd::RS<255, 255 - 2 * 17> coder17;
ezpwd::RS<255, 255 - 2 * 18> coder18;
ezpwd::RS<255, 255 - 2 * 19> coder19;
ezpwd::RS<255, 255 - 2 * 20> coder20;
ezpwd::RS<255, 255 - 2 * 21> coder21;
ezpwd::RS<255, 255 - 2 * 22> coder22;
ezpwd::RS<255, 255 - 2 * 23> coder23;
ezpwd::RS<255, 255 - 2 * 24> coder24;
ezpwd::RS<255, 255 - 2 * 25> coder25;
ezpwd::RS<255, 255 - 2 * 26> coder26;
ezpwd::RS<255, 255 - 2 * 27> coder27;
ezpwd::RS<255, 255 - 2 * 28> coder28;
ezpwd::RS<255, 255 - 2 * 29> coder29;
ezpwd::RS<255, 255 - 2 * 30> coder30;
ezpwd::RS<255, 255 - 2 * 31> coder31;
ezpwd::RS<255, 255 - 2 * 32> coder32;
ezpwd::RS<255, 255 - 2 * 33> coder33;

std::vector<uint8_t> encode(uint32_t k, std::vector<uint8_t> &data) {
    switch (k) {
        case 1:
            return encode<1>(coder1, data);
        case 2:
            return encode<2>(coder2, data);
        case 3:
            return encode<3>(coder3, data);
        case 4:
            return encode<4>(coder4, data);
        case 5:
            return encode<5>(coder5, data);
        case 6:
            return encode<6>(coder6, data);
        case 7:
            return encode<7>(coder7, data);
        case 8:
            return encode<8>(coder8, data);
        case 9:
            return encode<9>(coder9, data);
        case 10:
            return encode<10>(coder10, data);
        case 11:
            return encode<11>(coder11, data);
        case 12:
            return encode<12>(coder12, data);
        case 13:
            return encode<13>(coder13, data);
        case 14:
            return encode<14>(coder14, data);
        case 15:
            return encode<15>(coder15, data);
        case 16:
            return encode<16>(coder16, data);
        case 17:
            return encode<17>(coder17, data);
        case 18:
            return encode<18>(coder18, data);
        case 19:
            return encode<19>(coder19, data);
        case 20:
            return encode<20>(coder20, data);
        case 21:
            return encode<21>(coder21, data);
        case 22:
            return encode<22>(coder22, data);
        case 23:
            return encode<23>(coder23, data);
        case 24:
            return encode<24>(coder24, data);
        case 25:
            return encode<25>(coder25, data);
        case 26:
            return encode<26>(coder26, data);
        case 27:
            return encode<27>(coder27, data);
        case 28:
            return encode<28>(coder28, data);
        case 29:
            return encode<29>(coder29, data);
        case 30:
            return encode<30>(coder30, data);
        case 31:
            return encode<31>(coder31, data);
        case 32:
            return encode<32>(coder32, data);
        case 33:
            return encode<33>(coder33, data);
        default:
            return {};
    }
}

std::vector<uint8_t> decode(uint32_t k, std::vector<uint8_t> &data, std::vector<int>& missing) {
    switch (k) {
        case 1:
            return decode<1>(coder1, data, missing);
        case 2:
            return decode<2>(coder2, data, missing);
        case 3:
            return decode<3>(coder3, data, missing);
        case 4:
            return decode<4>(coder4, data, missing);
        case 5:
            return decode<5>(coder5, data, missing);
        case 6:
            return decode<6>(coder6, data, missing);
        case 7:
            return decode<7>(coder7, data, missing);
        case 8:
            return decode<8>(coder8, data, missing);
        case 9:
            return decode<9>(coder9, data, missing);
        case 10:
            return decode<10>(coder10, data, missing);
        case 11:
            return decode<11>(coder11, data, missing);
        case 12:
            return decode<12>(coder12, data, missing);
        case 13:
            return decode<13>(coder13, data, missing);
        case 14:
            return decode<14>(coder14, data, missing);
        case 15:
            return decode<15>(coder15, data, missing);
        case 16:
            return decode<16>(coder16, data, missing);
        case 17:
            return decode<17>(coder17, data, missing);
        case 18:
            return decode<18>(coder18, data, missing);
        case 19:
            return decode<19>(coder19, data, missing);
        case 20: return decode<20>(coder20, data, missing);
        case 21: return decode<21>(coder21, data, missing);
        case 22: return decode<22>(coder22, data, missing);
        case 23: return decode<23>(coder23, data, missing);
        case 24: return decode<24>(coder24, data, missing);
        case 25: return decode<25>(coder25, data, missing);
        case 26: return decode<26>(coder26, data, missing);
        case 27: return decode<27>(coder27, data, missing);
        case 28: return decode<28>(coder28, data, missing);
        case 29: return decode<29>(coder29, data, missing);
        case 30: return decode<30>(coder30, data, missing);
        case 31: return decode<31>(coder31, data, missing);
        case 32: return decode<32>(coder32, data, missing);
        case 33: return decode<33>(coder33, data, missing);

        default:
            return {};
    }
}

//int main() {
//    std::std::std::vector<uint8_t> data(128);
//
//    for (int i = 0; i < 128; ++i) {
//        data[i] = i;
//    }
//
//    auto encoded = encode(20, 6, data);
//
//    for (int i = 0; i < 6; ++i) {
//        for (int j = 0; j < encoded[0].size(); ++j) {
//            encoded[i][j] = 255;
//        }
//    }
//
//    auto decoded = decode(20, 6, encoded);
//
//    return 0;
//}
