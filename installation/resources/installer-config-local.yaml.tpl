apiVersion: v1
kind: ConfigMap
metadata:
  name: application-connector-certificate-overrides
  namespace: compass-installer
  labels:
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.connector.caKey: "LS0tLS1CRUdJTiBSU0EgUFJJVkFURSBLRVktLS0tLQpNSUlKS1FJQkFBS0NBZ0VBb1Q5d3FIT2RVRWo2K0ZtaGdsckVhc1BjaXd3MHBLNC8wN3A3NjFWVzdCY1ZOc1lyCnd5dDdQUHhxbWxJODdHNk9UYTM3SlFKT0g5cjFKeVV2VEl6MG1rdnFrME1aa0tCT3Y2U0xEeUNwcDhnNE1HekMKS0dyV05ZQkhGY3Iwd0tLbDVvVXRYbjl0MEc4U1c4bU1GMXFaT2k5aUxkRVl4T2RuR3Jhb1p6OUU2dExUUFNBQwp0Qkc5UFJ0Q3BvRVA5WXl2L1dKMDdnOUx4bUsvc1lVbXBKMHpFZ2xwOWJhVERONTVDcmhQNE5zcUtkdS9LR3dwCis4SVM3MXBncFhLaGFna3kzclhVVnluY2RtcEhLWDRuQTRHSHRLTFIyQTU5UHIrWlZZK2xQV2FwdE54d1hGekYKQ0FRUlJtNWJvUUhFSFNSaE1tcHIvemF6Z1pTMGhubTNJL21pRk5mODkvc0VJalpnQWI1eFp6VFNDMzVLOFpFdwprd0Q5L2F1RUZEd3BvL0RiS3JLM21JNzhwWWp3bE0rMHp3T2tkTmxUdVpRVjZWYVNCUnpMNC9pQXdLZ0NXSU0zCmRvVEZzNmJUdytVeVYrbG4vN2RwMUFlbGtublM1WmQwWG1SSHg1dWVyaFFBRkNjWXc5Sk9uQVk5NW8vUVBqeEUKaFZpL21LZHZnUDZVWDlPWitUdnVQUHpOeWYrTkp4OUVOT3RQNUxWTnBQdVJnMmtNeDZWeG52VHlKOC8wUjUxRgo5bmpndnF3MXQ0NHJnQnkvcTRVaWsrWHZBTDR1VjRpZzJwVStwRWJ5bWp6SjJlM2kzQUpKV0NOYWVHRmt3NmFBCmlOWE5maThUaTVGN214ODBLYnlLQWZWUzA0MmpFSG0wdi9jek03aXE1d1FaRUVJRlA2cWQ5ODNHWms4Q0F3RUEKQVFLQ0FnRUFqUDBpYlRmQjZrd1ZuUWNKOENlYkxGc2JRRDBvM29FNWI5RFR2ejQ4Sld3OWNVb3ZRNVNHU2huTwp3Q1o5L0tEaUxrdWNsNHgvY04wTGsvR3dmTGVXdkQ3NjJVNUhVU3pLRGtrNkNiMGVlb1RYbElmVDhIRVI0Vy9MCk4rUGd3M3F6b203NTczRnVQRnlSNmMyOWYwSUpUbFhWKzRlanA2OUplSk1UaGt0TTRDSDg3NnBJa3RnYjVnMHEKNXRsY2NmQlVoVElNV1lib1U0dE9YMUswS2lVRlhaVDdvQXZHWWU4NFdNWTFtYjhvQzdlSFdqblJMNzlPdlJnQgovMGZPbVI5MzZrR0VhNzQvZFE2U01GYU1tRVV1dWlQUFphR3RveXIyVUZpc081YkRkazkwczEydUxjY1lyOE9ZCnZKd0Z0UkYxSnhia1hSK2dMd0l1SXBMVUxsRjhoR2ZEb3QvWUZvZGJBM3A0dlM1MEc4c212YWlROVhnVXdGdFoKMFdpWU5iQ2JKdlY1d0w2aXlPcTRqakFGeExqU2lqc20ybFpVcTZiVDY5YVFpR0xmbXJDYnJDRnJFSkNpU1F3WgpGenVrMGlXYnM4ZXRySWFnNjdtY0kvTlF4RkVQbXB2ZEdQTllGZ0RwUDdEMnlhUnA4MzVjeWJVd04vcUE2RnBGCkJHSEQ0ZWl6OWd0WlhwT25NYzBibk9KcG5qMlNpaUUwaS95a0NuTlZvRm1SeXladVY0YVlQd29laWpOZ1NjZWIKQTY1TThHaWtUdmdiemFaL3YvOVdlWElURGJVWjcvQTBwYlRDTjZZZ0xscjBpQlJPMUFmeGhNNkVmems4QlFUSQpidWJ5U3ZwQWhiRk83aTZLVVAvZW03eVV0ckxBekdOR2dnTFJMVEpRWGdYUWNSeUdVWkVDZ2dFQkFNOUVhZW5nCm9peE1QaE5jTlN0dENEekhER2RDSFRWcDN2U0s0bkNHRFhGY3N2aytqMGRRM2pQYWx2V3dDMjBYOTZWQlR1bzAKTFhhcGgwZ0xHYloyS3JWQ0tPOWZrZXpQVWMwOWMzNmp1OExFZ3NBZ0dGOUpXeDVlMUVleUIzS0w0RnIreXVVVwpLbDc0azltYmZVK1dRcG1vc3hKdW41ZjBqbzhoclRucGtqR2FqUHhTL05GaE1YVHpGdjl6dEhVblZhTzJVYzI3ClB1WUJNSkxGdTQ2Y2xxWG53RTQxUTVwVHFLZGwxaUw0dFRrVER5ekFNUHZucW9HSld4Q1ZidjR3b2Z1WXBrY0gKSG50SlVEMVhhSnI1N0xoSU4zeDRCN2pFU3JYTUwrOVVIcWgzTDRna2FwV0VTTDlsY1d3aWxlcHhYdG5RdDI3NgprK0NXVVpuY09TUklVM2NDZ2dFQkFNY3BGRGdBWHVYS040b0Z3TC9lRkkvUWVSeXBDakQ2LytidVdCV2F0SXdQCnprQ1h6ODNaY256NmQ4YUlHc2I1bi9tbFYwaHBCdkdtZzBkN0ZHbHZoWDVhZDRzeXFLZWZJajZtWG1PMDRrQ3gKa29kZkErUVh5Q1EwWlFSRi9TenRyNWxWQi9yYnhhYzdCVlplQUxvVUtSSzdUMFhCOHczVXN1SnNyNGR0QVkwZQpDUDlmQlRybFowYjBFSnROWXlVbGFSNWRkaFVuYnJFeHFkUzg5Yk8ybUxBUjVDTVJJRXZScTRseGlQYzRTTXhYCmhUY2FSTmpBU2M3WE5NUUI1WVpXeUxaYzN1RndQaXRHK0sxTGMwNWs1VFpWeEtNdmVhOXBiZkRHLysyQm9IcHIKOEFHVXVKVlk2bCtGTzNleXlmenRPVGtpMUkwNmlxU2dURzVrbzJ0bXlla0NnZ0VBUjVYZWFzdU4xM1RodjdnSwpHUnlJU3MySXFDVTZoMWN3alE5bTAreEl1azJFOXZhM2I2OHJmNGRRdWp4NlJjeVFXTUFzckZFbkhxUEF1STQwCjdFTDF6ekt4aHJOZ2FBVFd3T2NuZTZhN1U3S2hZZy96dXYxUC9qWk1aUkxFNWJnUDNmM0FQODBmQnp3ZGZIdnEKbE5GVjRWSlZ2dGo4UC9SVVJIVWlLaTFVczlNb1BJSEJGZVBXdkFpMWViY1JyYURQUUVMWkVCQkswZys1SWdndgpGanRaQUtZQlVrR3RQcUVFVUFTcEo5ejBZbWtGeGJQL2R4RjFYMVg4WU1icjFka2dLUkI0NVhFOUF1RzRWK2RYCmxxY1pMakNyRVU4M2c0WXdNNGY1U2xTb1hoRUVGcVpWTlp6QnIzRXU4bVVqbUJ4ZDRTYm9JK2xocDZEalFCdkMKbEpoeVV3S0NBUUI4QUpyREo0L3VtVkt0VUZtcjNQV0dlYkgrNDAwaUpCWFRUbEZ2Mml4U0RNRkp2SHc1V2d1TAp2MU4yUEdZWHYzTVl1QmE1VWhOdHdGUjYzQ3BnWDN5SnFJQklIaG1lakZtQkVvc3duMzVEODR3ZFYwNlA1VExMClFBZ3BlZjVoeS9nS2kwUDFzSUxIVmR0RDVER2xxa25Nak8yVnJHWE9GY0h2Y3VaemRxNkJrOUxjVmVobXZGRHEKZjZvYldEckQ5U0FYTlBBQnlkU0U1VHd0NWgxQmNRNXVpaVUycEVJc2t2YXdGQTNJaDdYajdSWlhzYlp1RW9PaQpFcUthNitkaUZvVFA3dEVqSW9UQzQyU1FXYXNJZzQrbm5nMVo0WVJ0Y0VKd3FTYk9WV2g2OE51MTBFaUJUS1JaCkp4Wll0K3hGMjlwR05lYUxySWlJYWZwTXZjSjJhOENKQW9JQkFRRE5Nbk43M1k4UWNwSld4eTk4SmNLSklrTTQKY3BIeUM1U2dva2QvRWdBZzdMdW96a0grQWtWREZTUmcxdmRDSjdSMU9yTXNxeUVyYy9xV0pDYWpqVWJ6UlVyOApJQlN5b2Y1Zko3UHZ1VFc2WkV6RUR2WEw5WW90VExLM1VnUWpGc0N6c3dIUVhOaEE1QVBHWWZZa2hsTnlubVdpCnpyYlFYU1Q2RkxDeSt3UW9UQ3NzVHF4NW8yOEZiZ21RSTZvQzdmQ0ViTzZNL2dCTUxtbjhFNW03dWJVWE9wZmYKaEhsVGcwaVpKU0lDTDNRK25aZ0J4dVNSZ1hsMjJkUjZNYndiRDR6TTltNFhjTjNLTkNGMFY2bkJuQ0dnQ2hNaQp3aWpnS0w5RU1LemZFN0N4K05ZR0x5em9icVBXWDdMak81UjlPUWlpa3FBNmxxTldlN1g0b2NacFZ6c28KLS0tLS1FTkQgUlNBIFBSSVZBVEUgS0VZLS0tLS0K"
  global.connector.caCertificate: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUYwakNDQTdxZ0F3SUJBZ0lKQUtudzEwQi9zejJUTUEwR0NTcUdTSWIzRFFFQkN3VUFNRTB4Q3pBSkJnTlYKQkFZVEFrSkhNUTR3REFZRFZRUUlEQVZzYjJOaGJERU9NQXdHQTFVRUJ3d0ZiRzlqWVd3eERqQU1CZ05WQkFvTQpCV3h2WTJGc01RNHdEQVlEVlFRRERBVnNiMk5oYkRBZ0Z3MHlNakF5TWpVeE16VTFNVEJhR0E4eU1USXlNREl3Ck1URXpOVFV4TUZvd1RURUxNQWtHQTFVRUJoTUNRa2N4RGpBTUJnTlZCQWdNQld4dlkyRnNNUTR3REFZRFZRUUgKREFWc2IyTmhiREVPTUF3R0ExVUVDZ3dGYkc5allXd3hEakFNQmdOVkJBTU1CV3h2WTJGc01JSUNJakFOQmdrcQpoa2lHOXcwQkFRRUZBQU9DQWc4QU1JSUNDZ0tDQWdFQW9UOXdxSE9kVUVqNitGbWhnbHJFYXNQY2l3dzBwSzQvCjA3cDc2MVZXN0JjVk5zWXJ3eXQ3UFB4cW1sSTg3RzZPVGEzN0pRSk9IOXIxSnlVdlRJejBta3ZxazBNWmtLQk8KdjZTTER5Q3BwOGc0TUd6Q0tHcldOWUJIRmNyMHdLS2w1b1V0WG45dDBHOFNXOG1NRjFxWk9pOWlMZEVZeE9kbgpHcmFvWno5RTZ0TFRQU0FDdEJHOVBSdENwb0VQOVl5di9XSjA3ZzlMeG1LL3NZVW1wSjB6RWdscDliYVRETjU1CkNyaFA0TnNxS2R1L0tHd3ArOElTNzFwZ3BYS2hhZ2t5M3JYVVZ5bmNkbXBIS1g0bkE0R0h0S0xSMkE1OVByK1oKVlkrbFBXYXB0Tnh3WEZ6RkNBUVJSbTVib1FIRUhTUmhNbXByL3phemdaUzBobm0zSS9taUZOZjg5L3NFSWpaZwpBYjV4WnpUU0MzNUs4WkV3a3dEOS9hdUVGRHdwby9EYktySzNtSTc4cFlqd2xNKzB6d09rZE5sVHVaUVY2VmFTCkJSekw0L2lBd0tnQ1dJTTNkb1RGczZiVHcrVXlWK2xuLzdkcDFBZWxrbm5TNVpkMFhtUkh4NXVlcmhRQUZDY1kKdzlKT25BWTk1by9RUGp4RWhWaS9tS2R2Z1A2VVg5T1orVHZ1UFB6TnlmK05KeDlFTk90UDVMVk5wUHVSZzJrTQp4NlZ4bnZUeUo4LzBSNTFGOW5qZ3ZxdzF0NDRyZ0J5L3E0VWlrK1h2QUw0dVY0aWcycFUrcEVieW1qekoyZTNpCjNBSkpXQ05hZUdGa3c2YUFpTlhOZmk4VGk1RjdteDgwS2J5S0FmVlMwNDJqRUhtMHYvY3pNN2lxNXdRWkVFSUYKUDZxZDk4M0daazhDQXdFQUFhT0JzakNCcnpBUEJnTlZIUk1CQWY4RUJUQURBUUgvTUIwR0ExVWREZ1FXQkJUbwprMzBXVzhCRUxnMEZBWEZscjNnYnpkektvVEI5QmdOVkhTTUVkakIwZ0JUb2szMFdXOEJFTGcwRkFYRmxyM2diCnpkektvYUZScEU4d1RURUxNQWtHQTFVRUJoTUNRa2N4RGpBTUJnTlZCQWdNQld4dlkyRnNNUTR3REFZRFZRUUgKREFWc2IyTmhiREVPTUF3R0ExVUVDZ3dGYkc5allXd3hEakFNQmdOVkJBTU1CV3h2WTJGc2dna0FxZkRYUUgregpQWk13RFFZSktvWklodmNOQVFFTEJRQURnZ0lCQUpkbHNtR2k0d1hvSlZ4SnlKVzlDWEZPejhZWkhTWlhicEdsClJXcmI4QkZIMy9SNFhPTTQ5Y3Y4UzErYUZvQ3hJd0taRGFhZVIwNVIxK05jVXNPSnhZL2tGWXNHN3kvMTJFRVIKTi90anVRTGhPNnEzT1piZTUwUFErS2pxbmxURnQrT1ptV1ZQc3NzbUU2WWNhWXBOc09FcmVWZWxNWHNybDBndgpLZ3hLUXFJQ3hJTThNTG1QeUUxSURsS3RSL0RrTkNqNDQybmgvNVlwaUFoOG1BTUE5QjFtcE1uSTZpZFB0RzRhCm5uSGxiVlZaMUE4UWF2bXlCNGRueVZwNlB3QnhSMXJjN0xoV3VBT3V2WkNWWGpxUVZLMkREdEJ0Q3RGcUUrQlQKa1I4MnFjMGtKR2IwYTFkRWRubmJEdU1BQzA0WnRrQnpQTlV6RzNqSkNKVkZhVEsyMVhZQjNkMUR4UTVSUFJyTQptZmVXUlM5cytmdlplOTZmb1BIMmZuREN0aFJkOXI5UXdWSmhZMkRlTzlQOVd0TW9QMWhMempHSnQ1czB2RmpRClE2RDdZaXFHajkvTXhoc1BqZjlucjhCL3ZUVzFWV1hackh2NDBlR01tM1kwRXUvMlA0bFR5YnEyUlp2enhlQTgKTmlXL0lEa3VrbmlMaU95UnZKL1RWVXFsYllhbm9qeUVlbzJaTnp0RTVadTRYUndabnZIeVIzNFVrMldTYjdOSAp4Tk53TUt6ZHdROFpQQnM1OCtZUjByUGlGN2V4WXprVnFOdHNydEYwUjVPRm5kZjN4Yi9GeC8rRTBhd1UyUlZSClNiTzUxM1JIMTFZZ3RneWlMM0kxRjk3NGFDZTRhYkJuWHdmWDA1eUl2a2dEbFAvdXA5T25WRjV6eU8xSkRmUDAKSlFqY20zbGsKLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-certificate-overrides
  namespace: compass-installer
  labels:
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.ingress.tlsCrt: "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUU4RENDQXRpZ0F3SUJBZ0lKQUpBKzlrQzZZZnZlTUEwR0NTcUdTSWIzRFFFQkN3VUFNQmN4RlRBVEJnTlYKQkFNTURDb3VhM2x0WVM1c2IyTmhiREFlRncweE9EQTNNVGd3T0RNNE1UTmFGdzB5T0RBM01UVXdPRE00TVROYQpNQmN4RlRBVEJnTlZCQU1NRENvdWEzbHRZUzVzYjJOaGJEQ0NBaUl3RFFZSktvWklodmNOQVFFQkJRQURnZ0lQCkFEQ0NBZ29DZ2dJQkFPZ05XbVROOU1ranlrUjdvQ0JGSFVqL01EcWh5bml3NEJITGo4ZTdDTFV5dFdwVHZXTkoKU1FiaFZDK3c5NkhHbU50MHZGUTR4OExUa3NNUmorcVZrdkcwKzBDTE1WQm14UjBVdnpZR09QRVRKREtsNTZkcQppMEM1Y2dnU3dkNjcveWxRZmQxc3FEVVBHM0pVZlNOSFFRWSs1SkUwaUpjZXhCQ3cvS200UXlUQXB0aEwzdEgxCm1ZRFBBQ0hUdVpsbGc2RVN4Z2RKY20rVmg5UkRvbWlON250ZjBZVG1xV1VIZXhZUkUrcTFGY0VhUHg5L2Q5QUwKZWd1WXZHTkduMVA4K1F3ZE5DKzVaMEVGa1EyS3RtalRBR09La2xJUG5NQkZubXhGSEtNZzNTa3RkYzYzSU5TNwpjSDM3dDNjb2l4N09HcDllcnRvSFZ5K2U5KzdiYnI3Z0lhWHNORVZUcjRGWXJIMlNPeHI1MjVRQUpHTys1dllQCkNYSlJJdXFNZWU3ZTg3aDRIU3JFVVdWTnhOdElBR1ZNYlRtV2F6aDJsSEFoazJjZVIzUkRaQzB3NGFOd3NRU1UKZE1yMmFCaUFobmFNSG40YnZIZnd1OHFRcjF5aVRUNG00SkltRkZTNDN2bGNETDZmVUJnNnF0U2p6QjF0WTFjTwp2Mnl1QXQzR1lKSENVMFd4R1ZEQk80T0ZncEc3aE5SRHoyRzVFMXQvRWU2VDUxdys2K295eXhib01pK09kVWdNCnBZOGlzcDdBTWhtdlZFOUdtNUhwQkpBRjYvcmE1WUNnNDh1SWtybkdZbm85eVNKcGhLU0JsMFRML01PNVFRaUoKb0hFNnV5ZkJXQjlZZ1NBdDl0MjJqU3FSWkpSdEtrbHB6bFJIRkROWTRwZlRraXlPSlFXYlM4RXJBZ01CQUFHagpQekE5TUFrR0ExVWRFd1FDTUFBd0N3WURWUjBQQkFRREFnWGdNQ01HQTFVZEVRUWNNQnFDQ210NWJXRXViRzlqCllXeUNEQ291YTNsdFlTNXNiMk5oYkRBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQWdFQUlYYTlwenlyQ2dzMTRTOHUKZVFZdkorNEFzUE9uT1RGcExkaVl5UkVyNXdyNmJuMXUvMjZxc2FKckpxbkkyRk16SmdEQVRwZEtjbXRHYjBUOQp3S2wrYUJHcFFKcThrUWJwakVGTHhaWDJzaUNrRG82WittaUcrRjRKMHpKa3BKK0JHMS92eGZKbk0zK1ptdXQ5Ck9RV2ZjYTN3UHlhTWRDbGIyZjQwYlRFaFo5Mk9kcWlQMzFMbDlHWExSZmhaNTNsUzF0QWdvUGZoR25NbFY4b2MKWmxuSUROK25wS0Nma2tXUDJZUjlRLy9pa01tM01YRm9RSFppaVJseVZHSGFKZWRLMmNOQzlUYk4xNDFTaWZHZQo3V2FsQVBNcWNOQ3F3YStnN2RFSmR3ZjlRMklJTml0SjlDUVprT1dUZElYY2VHK2lZWWUrQXpmK1NkaHBocVdPCllFcDF6ek40dXI5U2VxU3NSaU9WY0RzVUFSa1M0clgrb0Vzb2hHL1Q5OTcrSDhjR2gzczl6TE84emtwRXZKSmEKS05QT0N5ODhVeEFOV2RRejFLMXRKVVQ2c3hkd0FEcXRJQnNPemhYVjlybDRRNStlZExlcmZPcUtCbUFRMUY5Swo2L1l0ZlNyY0JpeXZEU24wdFJ3OHJLRFVQU1hFNDFldXArOURNeThLVGl6T0RPTXVMSnR2dkJrTEFpNGNYQjVBCjQxMjBEdHdZQXNyNzNZYVl2SW8rWjV2OGZ4TjF3M3IwYS9KOVhZQlg3S3p1OFl4MnNUNWtWM2dNTHFCTXBaa3gKY29FTjNSandDMmV4VHl6dGc1ak1ZN2U4VFJ4OFFTeUxkK0pBd2t1Tm01NlNkcHFHNTE3cktJYkVMNDZzbkd0UgpCYUVOK01GeXNqdDU3ejhKQXJDMzhBMFN5dTQ9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0K"
  global.ingress.tlsKey: "LS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tCk1JSUpRd0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQ1Mwd2dna3BBZ0VBQW9JQ0FRRG9EVnBremZUSkk4cEUKZTZBZ1JSMUkvekE2b2NwNHNPQVJ5NC9IdXdpMU1yVnFVNzFqU1VrRzRWUXZzUGVoeHBqYmRMeFVPTWZDMDVMRApFWS9xbFpMeHRQdEFpekZRWnNVZEZMODJCamp4RXlReXBlZW5hb3RBdVhJSUVzSGV1LzhwVUgzZGJLZzFEeHR5ClZIMGpSMEVHUHVTUk5JaVhIc1FRc1B5cHVFTWt3S2JZUzk3UjlabUF6d0FoMDdtWlpZT2hFc1lIU1hKdmxZZlUKUTZKb2plNTdYOUdFNXFsbEIzc1dFUlBxdFJYQkdqOGZmM2ZRQzNvTG1MeGpScDlUL1BrTUhUUXZ1V2RCQlpFTgppclpvMHdCamlwSlNENXpBUlo1c1JSeWpJTjBwTFhYT3R5RFV1M0I5KzdkM0tJc2V6aHFmWHE3YUIxY3ZudmZ1CjIyNis0Q0dsN0RSRlU2K0JXS3g5a2pzYStkdVVBQ1JqdnViMkR3bHlVU0xxakhudTN2TzRlQjBxeEZGbFRjVGIKU0FCbFRHMDVsbXM0ZHBSd0laTm5Ia2QwUTJRdE1PR2pjTEVFbEhUSzltZ1lnSVoyakI1K0c3eDM4THZLa0s5YwpvazArSnVDU0poUlV1Tjc1WEF5K24xQVlPcXJVbzh3ZGJXTlhEcjlzcmdMZHhtQ1J3bE5Gc1JsUXdUdURoWUtSCnU0VFVRODlodVJOYmZ4SHVrK2RjUHV2cU1zc1c2REl2am5WSURLV1BJcktld0RJWnIxUlBScHVSNlFTUUJldjYKMnVXQW9PUExpSks1eG1KNlBja2lhWVNrZ1pkRXkvekR1VUVJaWFCeE9yc253VmdmV0lFZ0xmYmR0bzBxa1dTVQpiU3BKYWM1VVJ4UXpXT0tYMDVJc2ppVUZtMHZCS3dJREFRQUJBb0lDQUZDN3ZKUlB0M2QzVlRyb1MvaU9NemNmCldaYzhqT1hhbThwMUtRdlRQWjlWQ2hyNUVXNEdwRHFaa0tHYkR6eWdqTFBsZEZSVkFPTCtteFAwK3o0aFZlTjAKRk9vS3cxaDJ1T042UVdBNVgvdzNyYU5WWnpndThFM1BkeVhwNkx0bWFzcmo3elpuUkVwWmZESVZ4UWZPRllobgp2enZwckEvdnEwVW5YbkJwNUNwWVFIUUdTWHFBMlN3Z1dLcHNNQ2wzVVFsc0w2dC9XU29MT3h1VmdGNmg2clBQCnpXUlFuK1MvYW9wdDNLRU82WWVxYXdXNVltVG1hVXE1aytseU82S0w0OVhjSHpqdlowWU8rcjFjWWtRc0RQbVUKejMxdll4amQzOVZKWWtJNi85Y0Fzdmo5YTVXM3ROYVFDZStTRW56Z05oRDJieHo1NnRKdG0xTGwweXpqYTdEVgp0OW9aQVd6WURMOTVnMDFZWVlQMG51c3A4WWlNOXd0OUtLc2NkcUd3RFN5WHNrWFBrcWFLVjNVcHJYdmhFbFVaCkErMmtjcm9VaDNGbEV4bGpwUmtKVkhOZXJ3NHRLRi9oYTFWRjZPdE10eTVQcXV0N0dGQmIvamtWeUg5cnpueWUKTXQyTWVyTTVPazMwd1NuTThISUdTUXpxYlplekJEZlNaUzRzcWdZZnBIMlhtMEs0SjgrRUowQ2hhMXZVSmVNMAoyZ284d00vaHljdmtqTEgxSmM3OEhpaVBTQ01udkpHemUxc2tWdmtRRFhBSFdldzBTUHpUSTZHYjZCb0Y4aVNECm0wZjR2azNoV3NlUWZBaXVZSnlUeUZXNmRhOGE1K2lpSDN4cVRsUUN1MDN1Nmo0U0l4aThJZlNmd0YwQTBldVAKNGtzalZTZVZyT3ExUnlvNUtpR3hBb0lCQVFEOWZtYnl6aW9QdVhRYzl0QXBxMUpSMzErQzlCdFFzcDg5WkZkSQpQaU5xaTJ3NVlVcTA0OFM2Z3VBb3JGOHNObUI3QjhWa1JlclowQ3hub2NHY0tleWdTYWsvME5qVElndk5weGJwCnBGbkFnRjlmbW1oTEl2SlF2REo2Q0ZidDRCQlZIdkJEWlYyQnZqK0k3NUxkK01jN2RPVDdFek1FRjBXcUdzY2MKTUpyNjRXQi9UMkF5dWR1YXlRT2NobmJFQ25FUmdRcHFlbG54MVBraytqbGNvYUs5QjFYUStVOUgzOHppM0FYNApENUxMY0Nhem9YYWlvS0swckNlNE5Ga2hOVXd0TFV0QXhSTXk2aksyZUZudWczUFRlY3N1WktNMElITktqZ0dCCnpGanZVb2tMcFVFb3BJa3FHM09yc0xmanpHZW9jaVFPUXNEdzlUb3lXL0FSOFhmWkFvSUJBUURxV0s2TThVN3EKUXJPeTYzNnpEZlBaZ2ljeXlsUWVoOFdMclBlbW5NeWdQYzR4eWoybnMrTmVRSnNEUmtPT2tWY250SEFYaTcwWgoyT0NCV3dwZHJuTXlSc3RIMU05bjdFNU5TVHZlZDlkU3YxUzRBb3NzS1hDSmgyUHBjYjV0OE9nL3ZGTlNYUlUyClk2aUorWTdOcDBZNDNxSlJOVnlRemd3YmFzaEpiUVdkVFFoVEVhdEVRS2JsUlZSblhlQWRjOXlhNUpHbkRpaTkKbFQrRWEzdFpvN1dha05oeHJkYjVuTkZ3a0xoNEs3ZkFtT3ZzMVBMQWx0SUZqeURCeDEvY1ZHblpDUDBVQmJqZgpkU2FueXBBdVRuMzd1VUwrcXpPVDlYWEZENllGT0x0NWV4d3RxdnUwSzZCNjcvajFFTlJDRk45RnMzWlV5RFFXClZUaDcybFhWU3NLakFvSUJBUUR6ZE5pdXpTNDhWK0thaHJpNXJGNnRYeGkrRG0vRmV5ZlFzSFBiWUVKbmEyd1AKVjgrR0YxS3p4a28vQmYySjJ0ZWlrWDRVcGNtK1UxNnlVUG8vWDB4eFRRMk55cWpUYmRsa005dWZuVWJOeVB6UQpOdDEvZkJxNVMyWTNLWmREY25SOUsrK1k2dHQ1Wmh4akNhUkdKMDVCWGkwa3JmWExNZ2FvTG51WUtWNVBJUEdxCms3TlNSSW9UQ0llOVpxN2Q3U0ZXckZZeW1UdVZOUFByZlo1bHhwOGphTTRVbTd4MnpReGJ2UERHb3o1YXdHV0wKRThGNncwaEF1UzZValVJazBLbE9vamVxQnh3L1JBcGNrUTNlTXNXbEQwNENTb2tyNFJhWlBmVllrY2ZBWWNaWgpOdWR6ZjBKMC9GU0ZTbjN4L0RoNTROV2NGS1IxUnpBVGVaVUJ4cVZSQW9JQkFRRFlDNmZvVWpOQnJ2ckNFVzkrCkhYZlk1Ni9CbUZ4U3hUTHU0U2h6VncwakViZTltVWljQ2pDc1hQMUwySVJCdEdaWU9YWTVqdDlvSzlSV0RSdVMKWUZqZFdmemduU1lWRmZyZUw0emRQVGlxbGEvQjhMNWptVlNoeGNycmxheE02Uk1FWjFlZGtDa1ZPbTFQdmwzVAo1TW5OZGhySXFWeE1OMWxjRVdiU29vclJpUW9Lb3poMHRQSG9YckZBbG9BZVJ3bHpWeE9jb21ZVzJiaDBHUzdmCjVoaHZoZWUxYmVISnY3UXFoWkU3WUhxSU9iTVBaUWJqWEdnRkxmMnlDRitzM2Jtem1DRFJTN0V6ZVdxSXVDdVMKTlZUYU0rSzZyQlRoN0NLRjZUWlNqQW55SmZoRmRlT1ZKNzlNZDEzYWVJaG0zNTB6UWc3dWZKL2drdkorNUR2TApacC9uQW9JQkFBVlg4WHpFTzdMVk1sbENKZFVUdURTdXNPcDNrQVlFZ2dZNFFRM3FNTlcxRnl5WEM3WjBGOWJFCmtTSEhkalJtU2RUbFZueGN2UW1KTS9WL2tJanpNUHhFT3NCS1BVVkR3N3BhOHdiejlGcTRPOCtJb3lqN1ZXclcKMmExL1FNWXlzSGlpTlBzNy8vWUtvMy9rdkhCWUY2SnNkenkyQkVSTkQ0aTlVOWhDN0RqcGxKR3BSNktMTVBsegpNWFJ3VjVTM2V3cnBXZVcxQW5ONC9EKy9zUGlNQTNnS0swSlBFdGVBV1dndEZHTnNBSkJnaFBoUExxQi9CcDUvCkhOeC96M0w0MWtqRnpqOHNWaHMrVDRZYlhiaGF2R2xxc2h5ZldQbnRhV1VOMG15MjU1RFdQUDhWa24yeFNlV2kKVm1hVW5TSDBTZ2tlUENMRnlra25yQzgxU2pXZkRBMD0KLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLQo="
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: compass-installation-config-overrides
  namespace: compass-installer
  labels:
    component: compass
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.isLocalEnv: "true"
  global.isForTesting: "true"
  global.domainName: "kyma.local"
  global.adminPassword: ""
  global.minikubeIP: ""
  global.ingress.domainName: "kyma.local"
  global.externalServicesMock.enabled: "true"
  global.externalServicesMock.auditlog.applyMockConfiguration: "true"
  global.auditlog.skipSSLValidation: "true"
  global.systemFetcher.enabled: "true"
  global.systemFetcher.systemsAPIEndpoint: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/systemfetcher/systems"
  global.systemFetcher.systemsAPIFilterCriteria: "no"
  global.systemFetcher.systemToTemplateMappings: '[{"Name": "temp1", "SourceKey": ["prop"], "SourceValue": ["val1"] }, {"Name": "temp2", "SourceKey": ["prop"], "SourceValue": ["val2"] }]'
  global.oathkeeper.mutators.authenticationMappingServices.tenant-fetcher.authenticator.enabled: "true"
  global.oathkeeper.mutators.authenticationMappingServices.subscriber.authenticator.enabled: "true"
  global.oathkeeper.mutators.authenticationMappingServices.nsadapter.authenticator.enabled: "true"
  gateway.gateway.auditlog.enabled: "true"
  gateway.gateway.auditlog.authMode: "oauth-mtls"
  system-broker.http.client.skipSSLValidation: "true"
  connector.http.client.skipSSLValidation: "true"
  operations-controller.http.client.skipSSLValidation: "true"
  global.systemFetcher.http.client.skipSSLValidation: "true"
  global.ordAggregator.http.client.skipSSLValidation: "true"
  global.http.client.skipSSLValidation: "true"
  global.tests.http.client.skipSSLValidation: "true"
  global.externalCertConfiguration.secrets.externalCertSvcSecret.manage: "true"
  global.externalServicesMock.oauthSecret.manage: "true"
  global.externalServicesMock.regionInstancesCredentials.manage: "true"
  global.tests.basicCredentials.manage: "true"
  global.tests.ordService.subscriptionOauthSecret.manage: "true"
  global.hydrator.http.client.skipSSLValidation: "true"

  global.tenantFetchers.account-fetcher.cron.enabled: "false"
  global.tenantFetchers.account-fetcher.enabled: "true"
  global.tenantFetchers.account-fetcher.dbPool.maxOpenConnections: "1"
  global.tenantFetchers.account-fetcher.dbPool.maxIdleConnections: "1"
  global.tenantFetchers.account-fetcher.job.interval: "10s"
  global.tenantFetchers.account-fetcher.manageSecrets: "true"
  global.tenantFetchers.account-fetcher.secret.name: "compass-account-fetcher-secret"
  global.tenantFetchers.account-fetcher.secret.clientIdKey: "client-id"
  global.tenantFetchers.account-fetcher.secret.clientSecretKey: "client-secret"
  global.tenantFetchers.account-fetcher.secret.oauthMode: "standard"
  global.tenantFetchers.account-fetcher.secret.oauthUrlKey: "url"
  global.tenantFetchers.account-fetcher.oauth.client: "client_id"
  global.tenantFetchers.account-fetcher.oauth.secret: "client_secret"
  global.tenantFetchers.account-fetcher.oauth.tokenURL: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080"
  global.tenantFetchers.account-fetcher.oauth.tokenPath: "/secured/oauth/token"
  global.tenantFetchers.account-fetcher.providerName: "external-svc-mock"
  global.tenantFetchers.account-fetcher.schedule: "0 23 * * *"
  global.tenantFetchers.account-fetcher.kubernetes.configMapNamespace: "compass-system"
  global.tenantFetchers.account-fetcher.kubernetes.pollInterval: "2s"
  global.tenantFetchers.account-fetcher.kubernetes.pollTimeout: "1m"
  global.tenantFetchers.account-fetcher.kubernetes.timeout: "2m"
  global.tenantFetchers.account-fetcher.fieldMapping.idField: "guid"
  global.tenantFetchers.account-fetcher.fieldMapping.nameField: "displayName"
  global.tenantFetchers.account-fetcher.fieldMapping.customerIdField: "customerId"
  global.tenantFetchers.account-fetcher.fieldMapping.discriminatorField: ""
  global.tenantFetchers.account-fetcher.fieldMapping.discriminatorValue: ""
  global.tenantFetchers.account-fetcher.fieldMapping.totalPagesField: "totalPages"
  global.tenantFetchers.account-fetcher.fieldMapping.totalResultsField: "totalResults"
  global.tenantFetchers.account-fetcher.fieldMapping.tenantEventsField: "events"
  global.tenantFetchers.account-fetcher.fieldMapping.detailsField: "eventData"
  global.tenantFetchers.account-fetcher.fieldMapping.entityTypeField: "type"
  global.tenantFetchers.account-fetcher.queryMapping.pageNumField: "page"
  global.tenantFetchers.account-fetcher.queryMapping.pageSizeField: "resultsPerPage"
  global.tenantFetchers.account-fetcher.queryMapping.timestampField: "ts"
  global.tenantFetchers.account-fetcher.query.startPage: "1"
  global.tenantFetchers.account-fetcher.query.pageSize: "1000"
  global.tenantFetchers.account-fetcher.shouldSyncSubaccounts: "false"
  global.tenantFetchers.account-fetcher.accountRegion: "local"
  global.tenantFetchers.account-fetcher.endpoints.accountCreated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/global-account-create"
  global.tenantFetchers.account-fetcher.endpoints.accountDeleted: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/global-account-delete"
  global.tenantFetchers.account-fetcher.endpoints.accountUpdated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/global-account-update"

  global.tenantFetchers.subaccount-fetcher.cron.enabled: "false"
  global.tenantFetchers.subaccount-fetcher.enabled: "true"
  global.tenantFetchers.subaccount-fetcher.dbPool.maxOpenConnections: "1"
  global.tenantFetchers.subaccount-fetcher.dbPool.maxIdleConnections: "1"
  global.tenantFetchers.subaccount-fetcher.job.interval: "10s"
  global.tenantFetchers.subaccount-fetcher.manageSecrets: "true"
  global.tenantFetchers.subaccount-fetcher.secret.name: "compass-subaccount-fetcher-secret"
  global.tenantFetchers.subaccount-fetcher.secret.clientIdKey: "client-id"
  global.tenantFetchers.subaccount-fetcher.secret.clientSecretKey: "client-secret"
  global.tenantFetchers.subaccount-fetcher.secret.oauthMode: "oauth-mtls"
  global.tenantFetchers.subaccount-fetcher.secret.clientCertKey: "client-cert"
  global.tenantFetchers.subaccount-fetcher.secret.clientKeyKey: "client-key"
  global.tenantFetchers.subaccount-fetcher.secret.oauthUrlKey: "url"
  global.tenantFetchers.subaccount-fetcher.secret.skipSSLValidation: "true"
  global.tenantFetchers.subaccount-fetcher.oauth.client: "client_id"
  global.tenantFetchers.subaccount-fetcher.oauth.secret: "client_secret"
  global.tenantFetchers.subaccount-fetcher.oauth.tokenURL: '{{ printf "https://%s.%s" .Values.global.externalServicesMock.certSecuredHost .Values.global.ingress.domainName }}'
  global.tenantFetchers.subaccount-fetcher.oauth.tokenPath: "/cert/token"
  global.tenantFetchers.subaccount-fetcher.providerName: "subaccount-fetcher"
  global.tenantFetchers.subaccount-fetcher.schedule: "0 23 * * *"
  global.tenantFetchers.subaccount-fetcher.kubernetes.configMapNamespace: "compass-system"
  global.tenantFetchers.subaccount-fetcher.kubernetes.pollInterval: "2s"
  global.tenantFetchers.subaccount-fetcher.kubernetes.pollTimeout: "1m"
  global.tenantFetchers.subaccount-fetcher.kubernetes.timeout: "2m"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.idField: "guid"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.nameField: "displayName"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.customerIdField: "customerId"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.discriminatorField: ""
  global.tenantFetchers.subaccount-fetcher.fieldMapping.discriminatorValue: ""
  global.tenantFetchers.subaccount-fetcher.fieldMapping.totalPagesField: "totalPages"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.totalResultsField: "totalResults"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.tenantEventsField: "events"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.detailsField: "eventData"
  global.tenantFetchers.subaccount-fetcher.fieldMapping.entityTypeField: "type"
  global.tenantFetchers.subaccount-fetcher.queryMapping.pageNumField: "page"
  global.tenantFetchers.subaccount-fetcher.queryMapping.pageSizeField: "resultsPerPage"
  global.tenantFetchers.subaccount-fetcher.queryMapping.timestampField: "ts"
  global.tenantFetchers.subaccount-fetcher.query.startPage: "1"
  global.tenantFetchers.subaccount-fetcher.query.pageSize: "1000"
  global.tenantFetchers.subaccount-fetcher.shouldSyncSubaccounts: "true"
  global.tenantFetchers.subaccount-fetcher.accountRegion: "local"
  global.tenantFetchers.subaccount-fetcher.subaccountRegions: "test"
  global.tenantFetchers.subaccount-fetcher.endpoints.subaccountCreated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-create"
  global.tenantFetchers.subaccount-fetcher.endpoints.subaccountDeleted: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-delete"
  global.tenantFetchers.subaccount-fetcher.endpoints.subaccountUpdated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-update"
  global.tenantFetchers.subaccount-fetcher.endpoints.subaccountMoved: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-move"
  global.director.fetchTenantEndpoint: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenants/v1/fetch"
  global.tenantFetcher.fieldMapping.idField: "guid"
  global.tenantFetcher.fieldMapping.nameField: "displayName"
  global.tenantFetcher.fieldMapping.customerIdField: "customerId"
  global.tenantFetcher.fieldMapping.discriminatorField: ""
  global.tenantFetcher.fieldMapping.discriminatorValue: ""
  global.tenantFetcher.fieldMapping.totalPagesField: "totalPages"
  global.tenantFetcher.fieldMapping.totalResultsField: "totalResults"
  global.tenantFetcher.fieldMapping.tenantEventsField: "events"
  global.tenantFetcher.fieldMapping.detailsField: "eventData"
  global.tenantFetcher.fieldMapping.entityTypeField: "type"
  global.tenantFetcher.endpoints.subaccountCreated: "http://compass-external-services-mock.compass-system.svc.cluster.local:8080/tenant-fetcher/subaccount-create"
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: compass-installation-oidc-kubeconfig-service-overrides
  namespace: compass-installer
  labels:
    component: oidc-kubeconfig-service
    installer: overrides
    kyma-project.io/installation: ""
data:
  global.isLocalEnv: "true"
  global.domainName: "kyma.local"
  global.minikubeIP: ""
  global.ingress.domainName: "kyma.local"
  config.oidc.caFile: "/etc/dex-tls-cert/tls.crt"
