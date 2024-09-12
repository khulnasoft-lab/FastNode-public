#include "fastnode_smc.h"

#define IOSERVICE_SMC "AppleSMC"
#define IOSERVICE_MODEL "IOPlatformExpertDevice"

#define DATA_TYPE_SP78 "sp78"

typedef enum {
  fastnode_kSMCUserClientOpen = 0,
  fastnode_kSMCUserClientClose = 1,
  fastnode_kSMCHandleYPCEvent = 2,
  fastnode_kSMCReadKey = 5,
  fastnode_kSMCWriteKey = 6,
  fastnode_kSMCGetKeyCount = 7,
  fastnode_kSMCGetKeyFromIndex = 8,
  fastnode_kSMCGetKeyInfo = 9,
} fastnode_selector_t;

typedef struct {
  unsigned char major;
  unsigned char minor;
  unsigned char build;
  unsigned char reserved;
  unsigned short release;
} fastnode_SMCVersion;

typedef struct {
  uint16_t version;
  uint16_t length;
  uint32_t cpuPLimit;
  uint32_t gpuPLimit;
  uint32_t memPLimit;
} fastnode_SMCPLimitData;

typedef struct {
  IOByteCount data_size;
  uint32_t data_type;
  uint8_t data_attributes;
} fastnode_SMCKeyInfoData;

typedef struct {
  uint32_t key;
  fastnode_SMCVersion vers;
  fastnode_SMCPLimitData p_limit_data;
  fastnode_SMCKeyInfoData key_info;
  uint8_t result;
  uint8_t status;
  uint8_t data8;
  uint32_t data32;
  uint8_t bytes[32];
} fastnode_SMCParamStruct;

typedef enum {
  fastnode_kSMCSuccess = 0,
  fastnode_kSMCError = 1,
  fastnode_kSMCKeyNotFound = 0x84,
} fastnode_kSMC_t;

typedef struct {
  uint8_t data[32];
  uint32_t data_type;
  uint32_t data_size;
  fastnode_kSMC_t kSMC;
} fastnode_smc_return_t;

static const int fastnode_SMC_KEY_SIZE = 4; // number of characters in an SMC key.
static io_connect_t fastnode_conn;          // our connection to the SMC.

kern_return_t fastnode_open_smc(void) {
  kern_return_t result;
  io_service_t service;

  service = IOServiceGetMatchingService(kIOMasterPortDefault,
                                        IOServiceMatching(IOSERVICE_SMC));
  if (service == 0) {
    // Note: IOServiceMatching documents 0 on failure
    printf("ERROR: %s NOT FOUND\n", IOSERVICE_SMC);
    return kIOReturnError;
  }

  result = IOServiceOpen(service, mach_task_self(), 0, &fastnode_conn);
  IOObjectRelease(service);

  return result;
}

kern_return_t fastnode_close_smc(void) { return IOServiceClose(fastnode_conn); }

static uint32_t fastnode_to_uint32(char *key) {
  uint32_t ans = 0;
  uint32_t shift = 24;

  if (strlen(key) != fastnode_SMC_KEY_SIZE) {
    return 0;
  }

  for (int i = 0; i < fastnode_SMC_KEY_SIZE; i++) {
    ans += key[i] << shift;
    shift -= 8;
  }

  return ans;
}

static kern_return_t fastnode_call_smc(fastnode_SMCParamStruct *input, fastnode_SMCParamStruct *output) {
  kern_return_t result;
  size_t input_cnt = sizeof(fastnode_SMCParamStruct);
  size_t output_cnt = sizeof(fastnode_SMCParamStruct);

  result = IOConnectCallStructMethod(fastnode_conn, fastnode_kSMCHandleYPCEvent, input, input_cnt,
                                     output, &output_cnt);

  if (result != kIOReturnSuccess) {
    result = err_get_code(result);
  }
  return result;
}

static kern_return_t fastnode_read_smc(char *key, fastnode_smc_return_t *result_smc) {
  kern_return_t result;
  fastnode_SMCParamStruct input;
  fastnode_SMCParamStruct output;

  memset(&input, 0, sizeof(fastnode_SMCParamStruct));
  memset(&output, 0, sizeof(fastnode_SMCParamStruct));
  memset(result_smc, 0, sizeof(fastnode_smc_return_t));

  input.key = fastnode_to_uint32(key);
  input.data8 = fastnode_kSMCGetKeyInfo;

  result = fastnode_call_smc(&input, &output);
  result_smc->kSMC = output.result;

  if (result != kIOReturnSuccess || output.result != fastnode_kSMCSuccess) {
    return result;
  }

  result_smc->data_size = output.key_info.data_size;
  result_smc->data_type = output.key_info.data_type;

  input.key_info.data_size = output.key_info.data_size;
  input.data8 = fastnode_kSMCReadKey;

  result = fastnode_call_smc(&input, &output);
  result_smc->kSMC = output.result;

  if (result != kIOReturnSuccess || output.result != fastnode_kSMCSuccess) {
    return result;
  }

  memcpy(result_smc->data, output.bytes, sizeof(output.bytes));

  return result;
}
