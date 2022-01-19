#include <stdlib.h>
#include <stdio.h>

int main() {

	int nbDunnoTime = 14322029;
	unsigned char* pDunnoTime = &nbDunnoTime;	

	int nbDunnoSpace = 74342035;
	unsigned char* pDunnoSpace = &nbDunnoSpace;	
	
	int nbDunnoTot = nbDunnoTime+nbDunnoSpace;
	unsigned char* pDunnoTot = &nbDunnoTot;

	unsigned char header[30] = {0};
	
	for (int i = 0 ; i < 4 ; i += 1) {
		header[i] = *(pDunnoTime+4-i-1);
		header[i+4] = *(pDunnoSpace+4-i-1);
		header[i+8] = *(pDunnoTot+4-i-1);	
	}
	
	fwrite(&header, 1, sizeof(header), stdout);


	return 0;
}
