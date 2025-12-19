The DNS header is exactly
12 bytes long and is the same format for both queries and responses. Below are examples of DNS header byte representations, particularly the first 4 bytes, with a detailed breakdown of their meaning. 
DNS Header Structure (12 Bytes Total)
The 12-byte header is structured as a continuous series of bytes, with fields of specific lengths (bits or bytes) as defined in RFC1035. The number fields are encoded in network byte order (big-endian). 
Byte Range 
	Field Name	Length	Description
Bytes 0-1	ID	2 bytes	A 16-bit identifier assigned by the program that generates the query. The same ID is copied into the corresponding reply.
Bytes 2-3	Flags	2 bytes	A 16-bit field containing various flags (QR, Opcode, AA, TC, RD, RA, Z, RCODE).
Bytes 4-5	QDCOUNT	2 bytes	The number of entries in the question section.
Bytes 6-7	ANCOUNT	2 bytes	The number of resource records in the answer section.
Bytes 8-9	NSCOUNT	2 bytes	The number of name server resource records in the authority section.
Bytes 10-11	ARCOUNT	2 bytes	The number of resource records in the additional records section.
Byte Representation Examples
Consider a standard DNS query for an 'A' record (IPv4 address) for "google.com" with recursion desired.
Example 1: DNS Query Header (First 4 bytes)
A typical query packet might start with the following hexadecimal bytes:
86 2A 01 20 
Here is the breakdown of the first four bytes:

    Bytes 0-1 (ID): 86 2A (or 0x862A)
        This is a unique, random 16-bit transaction ID (34346 in decimal) used to match the response to the query.
    Bytes 2-3 (Flags): 01 20 (or 0x0120)
        Converting 0x0120 to binary gives 00000001 00100000. This 16-bit value is composed of specific flags:
            QR (Query/Response, first bit): 0 (Query).
            Opcode (4 bits): 0000 (Standard query).
            AA, TC (2 bits): 00 (Not authoritative, not truncated).
            RD (Recursion Desired, 7th bit of first byte): 1 (Recursion is desired).
            RA, Z, RCODE (remaining bits): 0... (Recursion not available in query, reserved bits zeroed, no error code). 

Example 2: Corresponding DNS Response Header (First 4 bytes)
The corresponding response from the server for the same query might have the following bytes:
86 2A 81 80 

    Bytes 0-1 (ID): 86 2A
        The ID is identical to the query, confirming it is the matching response.
    Bytes 2-3 (Flags): 81 80 (or 0x8180)
        Converting 0x8180 to binary gives 10000001 10000000.
            QR: 1 (Response).
            Opcode: 0000 (Standard query response).
            AA, TC: 00.
            RD: 1 (Recursion was desired in the original query).
            RA (Recursion Available, 7th bit of second byte): 1 (Server supports recursion).
            Z, RCODE: 0... (No error code). 

Domain Name Representation in Bytes
Following the 12-byte header, the question section encodes the domain name in a specific length-value pair format. For example, querying "microsoft.com" is encoded in bytes as: 
09 m i c r o s o f t 03 c o m 00

    09: Length of the first label (9 bytes for "microsoft").
    03: Length of the second label (3 bytes for "com").
    00: A null byte indicating the end of the domain name. 