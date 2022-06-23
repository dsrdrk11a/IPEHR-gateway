# Milestone Delivery :mailbox:

* **Application Document:** [ipEHR application](https://github.com/filecoin-project/devgrants/issues/418)
* **Milestone Number:** 1

**Context** (optional)
In this milestone we've designed and developed an MH-ORM and medical data storage structure.

**Deliverables**

| Number | Deliverable | Link | Notes |
| ------------- | ------------- | ------------- |------------- |
| :heavy_check_mark: | License | [LICENSE](https://github.com/bsn-si/IPEHR-gateway/blob/develop/LICENSE) | Apache 2.0 license |
| :heavy_check_mark: | Testing Guide | [Readme.md](https://github.com/bsn-si/IPEHR-gateway/blob/develop/README.md#how-to) | "How To" guide |
| 1. :heavy_check_mark: | Design of the medical data storage | [See IPEHR-gateway/Storage](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_1/1_Storage/README.md) | This project will use Filecoin as the EHR document repository. The implementation of saving documents in Filecoin will be implemented in the next milestone. | 
| 2. :heavy_check_mark: | Design of the index storage | [See IPEHR-gateway/Index_design](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/2_Index_design#readme) | To store the indexes, it is planned to create a special smart contract on the blockchain. This will ensure security, data immutability and fast access. Each document that is added to the system is encrypted using the ChaCha20-Poly1305 algorithm. A separate key is generated for each document. | 
| 3. :heavy_check_mark: | Data encryption methods | [See IPEHR-gateway/Encryption](https://github.com/bsn-si/IPEHR-gateway/blob/develop/progress/Milestone_1/3_Encryption/README.md) | When stored in the repository, EHR documents are pre-encrypted using the ChaCha20-Poly1305 streaming algorithm with message authentication. To encrypt each document a unique key is generated - a random sequence of 256 bits (32 bytes) + a unique 96 bits (12 bytes). Document ID is used as an authentication tag. |
| 4. :heavy_check_mark: | Filtering functionality | [See IPEHR-gateway/Filtering](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/4_Filtering) | When new EHR documents are created, the homomorphically encrypted data they contain is placed in a special DataSearch index tree structure, so that later selections can be made from this data using AQL queries. |
| 5. :heavy_check_mark: | Record creation and update functionality | [See IPEHR-gateway/EHR_creation](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/5_EHR_creation_and_update) | IPEHR-gateway implements functions for creating and updating documents according to openEHR standards. The EHR Information Model version 1.1.0, the latest stable version at the time of development, was used. | 
| 6. :heavy_check_mark: | API | [See IPEHR-gateway/API](https://github.com/bsn-si/IPEHR-gateway/tree/develop/progress/Milestone_1/6_API) | The minimum basic version of the REST API has been implemented to support document handling according to the latest stable version of the openEHR specification. In the next milestones the API will be supplemented with other methods to fully comply with the openEHR specifications. | 
| 7. :heavy_check_mark: | Extensive testing | [Readme.md]() | Milestone deliverables testing video guide |
