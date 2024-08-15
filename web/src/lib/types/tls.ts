export interface Domain {
	main?: string;
	sans?: string[];
}

export interface ClientTLS {
	ca?: string;
	caOptional?: boolean;
	cert?: string;
	key?: string;
	insecureSkipVerify?: boolean;
}

export interface CertAndStores {
	stores?: string[];
}

export interface ClientAuth {
	caFiles?: string[];
	clientAuthType?: string;
}

export interface GeneratedCert {
	resolver?: string;
	domain?: Domain;
}

export interface Options {
	minVersion?: string;
	maxVersion?: string;
	cipherSuites?: string[];
	curvePreferences?: string[];
	clientAuth?: ClientAuth;
	sniStrict?: boolean;
	preferServerCipherSuites?: boolean;
	alpnProtocols?: string[];
}

export interface Store {
	defaultCertificate?: Certificate;
	defaultGeneratedCert?: GeneratedCert;
}

export interface Certificate {
	certFile?: string;
	keyFile?: string;
}
